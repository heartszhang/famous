using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Interactivity;
using famousfront.utils;
using famousfront.messages;
using GalaSoft.MvvmLight.Messaging;
using GalaSoft.MvvmLight.Threading;
namespace famousfront.controls
{
  /*
   * mouse-enter
   * mouse-stay if mouse-enter for 300ms
   * mouse-leave 
   * mouse-nocare if mouse-leave for 300ms
   */
  class ImageTipProviderBehavior : Behavior<UIElement>
  {
    public static readonly DependencyProperty ImageUrlProperty =
        DependencyProperty.RegisterAttached("ImageUrl", typeof(string), typeof(ImageTipProviderBehavior));


    public static string GetImageUrl(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return obj.GetValue(ImageUrlProperty) as string;
    }


    public static void SetImageUrl(DependencyObject obj, string value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(ImageUrlProperty, value);
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
      var iu = GetImageUrl(AssociatedObject);
      if (string.IsNullOrEmpty(iu))
        return;
      var prevaid = ++action_id;
      await Task.Delay(ServiceLocator.Flags.ImageTipHideDelay).ConfigureAwait(false);
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        if (AssociatedObject.IsMouseOver)
          return;
        if (action_id != prevaid)
          return;
        Messenger.Default.Send(new ImageTipRequest { image_uri = iu, open = false });
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }

    async void ImageMouseEnter(object sender, System.Windows.Input.MouseEventArgs e)
    {
      var iu = GetImageUrl(AssociatedObject);
      if (string.IsNullOrEmpty(iu))
        return;
      var prevaid = ++action_id;
      await Task.Delay(ServiceLocator.Flags.ImageTipShowDelay).ConfigureAwait(false);
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        if (prevaid != action_id)
          return;
        if (AssociatedObject.IsMouseOver)
          Messenger.Default.Send(new ImageTipRequest { image_uri = iu, open = true });
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
    protected override void OnDetaching()
    {
      base.OnDetaching();
    }
  }

  class ImageTipServiceHostBehavior : Behavior<UIElement>
  {
    public static readonly DependencyProperty ImageUrlProperty =
        DependencyProperty.RegisterAttached("ImageUrl", typeof(string), typeof(ImageTipServiceHostBehavior),
        new FrameworkPropertyMetadata(null,FrameworkPropertyMetadataOptions.Inherits));


    public static string GetImageUrl(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return obj.GetValue(ImageUrlProperty) as string;
    }


    public static void SetImageUrl(DependencyObject obj, string value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(ImageUrlProperty, value);
    }
    protected override void OnAttached()
    {
      Messenger.Default.Register<ImageTipRequest>(this, OnImageTipRequest);
      base.OnAttached();
      
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
        var iu = GetImageUrl(AssociatedObject);
        if (msg.open)
        {
          SetImageUrl(AssociatedObject, msg.image_uri);
          return;
        }
        if (AssociatedObject.IsMouseOver)
        {
          return;
        }
        SetImageUrl(AssociatedObject, null);
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
  }
}
