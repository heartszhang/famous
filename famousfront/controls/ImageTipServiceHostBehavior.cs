using System;
using System.Diagnostics;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Input;
using System.Windows.Interactivity;
using famousfront.messages;
using famousfront.utils;
using GalaSoft.MvvmLight.Messaging;
using GalaSoft.MvvmLight.Threading;

namespace famousfront.controls
{
  class ImageTipServiceHostBehavior : Behavior<UIElement>
  {
    public static readonly DependencyProperty FeedImageProperty =
      DependencyProperty.RegisterAttached("FeedImage", typeof(datamodels.FeedMedia), typeof(ImageTipServiceHostBehavior),
        new FrameworkPropertyMetadata(null,FrameworkPropertyMetadataOptions.Inherits));


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
    protected override void OnAttached()
    {
      Messenger.Default.Register<ImageTipRequest>(this, OnImageTipRequest);
      base.OnAttached();
      AssociatedObject.MouseLeave += (new MouseEventHandler(OnImageTipMouseUp)).MakeWeakSpecial(eh => AssociatedObject.MouseLeave -= eh);
    }
    int action_id;
    private async void OnImageTipMouseUp(object sender, MouseEventArgs e)
    {
      var prev = ++action_id;
      await Task.Delay(1200);
      if (prev != action_id)
        return;
      SetFeedImage(AssociatedObject, null);
    }
    protected override void OnDetaching()
    {
      Messenger.Default.Unregister(this);
      base.OnDetaching();
    }
    async void OnImageTipRequest(ImageTipRequest msg)
    {
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        ++action_id;
        var iu = GetFeedImage(AssociatedObject);
        if (msg.open)
        {
          SetFeedImage(AssociatedObject, msg.image);
          return;
        }
        if (AssociatedObject.IsMouseOver)
        {
          return;
        }
        if (iu == msg.image)
          SetFeedImage(AssociatedObject, null);
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
  }
}