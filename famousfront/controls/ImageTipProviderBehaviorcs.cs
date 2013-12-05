using System;
using System.Diagnostics;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Input;
using System.Windows.Interactivity;
using famousfront.utils;
using famousfront.messages;
using GalaSoft.MvvmLight.Messaging;
using GalaSoft.MvvmLight.Threading;
namespace famousfront.controls
{
  class ImageTipProviderBehavior : Behavior<UIElement>
  {
    public static readonly DependencyProperty FeedImageProperty =
        DependencyProperty.RegisterAttached("FeedImage", typeof(datamodels.FeedMedia), typeof(ImageTipProviderBehavior));


    public static datamodels.FeedMedia GetFeedImage(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return obj.GetValue(FeedImageProperty) as datamodels.FeedMedia;
    }


    public static void SetFeedImage(DependencyObject obj, datamodels.FeedMedia value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(FeedImageProperty, value);
    }

    int action_id;
    protected override void OnAttached()
    {
      base.OnAttached();
      AssociatedObject.MouseEnter += (new MouseEventHandler(ImageMouseEnter)).MakeWeakSpecial(eh => AssociatedObject.MouseEnter -= eh);
      AssociatedObject.MouseLeave += (new MouseEventHandler(ImageMouseLeave)).MakeWeakSpecial(eh => AssociatedObject.MouseLeave -= eh);
    }

    async void ImageMouseLeave(object sender, System.Windows.Input.MouseEventArgs e)
    {
      var iu = GetFeedImage(AssociatedObject);
      if (iu == null )
        return;
      var prevaid = ++action_id;
      await Task.Delay(ServiceLocator.Flags.ImageTipHideDelay).ConfigureAwait(false);
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        if (AssociatedObject.IsMouseOver)
          return;
        if (action_id != prevaid)
          return;
        Messenger.Default.Send(new ImageTipRequest { image = iu, open = false });
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }

    async void ImageMouseEnter(object sender, System.Windows.Input.MouseEventArgs e)
    {
      var iu = GetFeedImage(AssociatedObject);
      if (iu == null )
        return;
      var prevaid = ++action_id;
      await Task.Delay(ServiceLocator.Flags.ImageTipShowDelay).ConfigureAwait(false);
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        if (prevaid != action_id)
          return;
        if (AssociatedObject.IsMouseOver)
          Messenger.Default.Send(new ImageTipRequest { image = iu, open = true });
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
  }
}
